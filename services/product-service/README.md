# Product Service

Quản lý danh mục linh kiện PC: model, thông số kỹ thuật, giá trị chuẩn hóa (lookup)
và API tìm kiếm sản phẩm cho website. Go + Gin + PostgreSQL (pgx).

## Cấu trúc

```text
cmd/api/                 entrypoint: load config, connect DB, graceful shutdown
config/                  đọc biến môi trường
migrations/              golang-migrate, version dạng YYYYMMDDhhmmss
internal/
  app/                   wiring module + đăng ký route (/api/v1, /api/v1/admin)
  platform/              hạ tầng dùng chung, không chứa nghiệp vụ
    apperror/            ErrNotFound / ErrConflict / ErrInvalid / ErrMalformed / ErrTooLarge
    database/            pgxpool, RunInTransaction (transaction qua context), map lỗi Postgres,
                         UpdateBuilder (UPDATE từng phần)
    httpx/               response.go: error -> HTTP status; request_body.go: decode JSON chặt;
                         validation.go: Validator; params.go: path/query param; phân trang
    optional/            Field[T]: phân biệt field bị bỏ qua và null trong request
  modules/
    catalog/             mỗi entity: dữ liệu, request, handler, service, repository
                         (component_type, manufacturer, product_model)
    lookup/              handler -> service (registry) -> repository dùng chung cho mọi bảng
    specs/               handler -> service (registry, transaction) -> Component
      cpu/               spec.go, patch.go, request.go, component.go, service.go, repository.go
    discovery/           handler -> service -> repository (+ dependencies.go)
```

Mỗi module tự đăng ký route qua `RegisterRoutes(admin, public)`.

### Trách nhiệm và cách đặt tên

- `module.go` chỉ khởi tạo dependencies và đăng ký route; không xử lý HTTP
  hoặc kiểm tra nghiệp vụ.
- `*_handler.go` / `handler.go`: đọc request, decode + validate request, gọi service và trả HTTP.
- `*_service.go` / `service.go`: kiểm tra nghiệp vụ, chọn loại dữ liệu, điều phối
  repository và transaction. Không parse JSON, không phụ thuộc Gin.
- `*_repository.go` / `repository.go`: thực thi SQL và ánh xạ kết quả.
- `*_request.go` / `request.go`: struct JSON request, hàm `Validate()` (kiểm tra
  từng field, viết tường minh bằng `httpx.Validator`, không dùng `binding` tag) và
  chuyển đổi sang input có kiểu cho service. Ví dụ `CreateProductModelRequest.Input()`.
  Mọi request phải cài `Validate()` (interface `httpx.Request`), nên quên validate
  sẽ báo lỗi lúc compile.
- File dữ liệu đặt theo đúng thực thể: `product_model.go` chứa `ProductModel`
  (model linh kiện, tương ứng bảng `models`), `cpu/spec.go` chứa `Spec` và
  `Attributes`, `cpu/patch.go` chứa `Patch`, `lookup/item.go` chứa `Item`.
  Không dùng tên chung `model.go`. Dữ liệu trả về có thể giữ `json` tags để bảo
  toàn response. Chỉ tách response struct riêng khi cần một hình dạng khác.
- `dependencies.go` (discovery) khai báo các interface mà module cần từ bên ngoài.
- `component.go` (specs) khai báo interface `Component`; `cpu/component.go` cài nó.

Luồng ghi CPU:

```text
Specs Handler -> Specs Service.ComponentFor(model) -> cpu.Component.DecodeRequest (request.go) -> Patch
              -> Specs Service.Upsert (transaction, kiểm tra loại model)
              -> CPU Service (gộp Patch vào spec hiện có, rule CPU) -> CPU Repository
```

CPU không có handler HTTP riêng vì dùng chung endpoint Specs. Mỗi loại linh kiện
cài interface `specs.Component` (giống abstract base class trong C++), gồm cả
`DecodeRequest`; chỉ handler gọi `DecodeRequest`. Service không nhận `[]byte`,
còn CPU service nhận `Patch` có kiểu cụ thể, khóa dòng spec hiện có
(`SELECT ... FOR UPDATE`), gộp patch rồi kiểm tra rule giữa nhiều field trên kết
quả đã gộp.

Code dùng chung cho nhiều bảng chỉ có ở `lookup`, vì các bảng lookup chỉ là danh
sách giá trị và không có nghiệp vụ riêng. Các entity nghiệp vụ như `component_types` và
`manufacturers` được viết tường minh, dù cấu trúc bảng hiện tại giống lookup, để
có thể thêm trường hoặc rule riêng mà không ảnh hưởng chỗ khác.

Quy tắc phụ thuộc:

- `catalog`, `lookup` không phụ thuộc module khác.
- `specs` chỉ dùng `catalog` qua interface `ModelTypeResolver`.
- `discovery` tự query DB bằng `Repository`; specs và SKU đi qua interface
  (`SpecReader`, `OfferReader` trong `dependencies.go`) được nối trong `internal/app`.

## API

| Nhóm | Endpoint | Ghi chú |
|---|---|---|
| Catalog (admin) | `/api/v1/admin/component-types`, `/manufacturers` | CRUD |
| | `/api/v1/admin/models` | CRUD, filter `q`, `component_type_id`, `manufacturer_id`, `is_active` |
| Specs (admin) | `GET/PUT/DELETE /api/v1/admin/models/:id/specs` | payload theo loại linh kiện của model |
| Lookup (admin) | `GET /api/v1/admin/lookups` | danh sách loại lookup |
| | `/api/v1/admin/lookups/:type[/:id]` | CRUD, ví dụ `:type = cpu-sockets` |
| Lookup (public) | `GET /api/v1/lookups?types=cpu-sockets,memory-types` | giá trị active, gom nhiều loại trong 1 request |
| | `GET /api/v1/lookups/:type` | |
| Discovery (public) | `GET /api/v1/products` | `q`, `component_type`, `manufacturer`, `sort=name\|newest`, `page`, `page_size` |
| | `GET /api/v1/products/:id` | model + specs + offers |

Quy ước:

- Danh sách trả `{items, page, page_size, total}`. `page_size` tối đa 100.
- Lỗi trả `{"error": "..."}` với:
  - 400: body không parse được (JSON sai, field lạ, sai kiểu, không phải một object).
  - 404; 409 (trùng hoặc đang được tham chiếu); 413 (body > 1 MB).
  - 422: vi phạm validate/nghiệp vụ, FK không tồn tại, giá trị vượt giới hạn cột.
- Mọi `PUT` cập nhật là cập nhật từng phần: field không gửi thì giữ nguyên, chỉ
  đặt `null` khi request gửi `null` rõ ràng. Field bắt buộc (`code`, `name`,
  `model_name`, `manufacturer_id`, `is_active`) không nhận `null` (422).
  `PUT /models/:id/specs` cũng vậy; lần ghi đầu tiên các field không gửi là `null`.
- Không xóa được lookup đang được dùng (409); hãy đặt `is_active=false`.
- `code` và `name` của component type không đổi được sau khi tạo (422), vì
  module spec và client nhận diện loại linh kiện qua chúng.
- `component_type_id` của model không đổi được sau khi tạo, vì nó quyết định
  bảng spec. `manufacturer_id` là bắt buộc.
- Admin API thấy cả dữ liệu inactive. Public API không bao giờ trả dữ liệu
  inactive: model/component type/manufacturer inactive bị ẩn khỏi
  `/products`, và tham chiếu tới lookup inactive trong specs bị bỏ
  (`null` hoặc loại khỏi mảng).

## Thêm loại spec mới (ví dụ GPU)

1. Tạo `internal/modules/specs/gpu/` theo mẫu `cpu/` gồm `spec.go`, `patch.go`,
   `request.go`, `component.go`, `service.go` và `repository.go`. Request dùng
   `optional.Field`, có `Validate()` và chuyển sang `Patch`; service gộp patch vào
   spec hiện có và kiểm tra rule giữa nhiều field; repository ghi bảng chính và
   thay toàn bộ bảng con trong transaction do `specs.Service` mở. `GetPublic`
   phải ẩn tham chiếu tới lookup inactive.
2. `gpu.Component` cài interface `specs.Component` (thêm dòng
   `var _ specs.Component = (*Component)(nil)` để compiler kiểm tra), rồi đăng ký
   `gpu.NewComponent(db)` trong `internal/app/app.go`.
3. `ComponentTypeCode()` của component và của `Patch` phải khớp
   `component_types.code` (ví dụ `gpu`).

CPU request giữ nguyên JSON phẳng. Decoder chỉ chấp nhận đúng một JSON object;
`null`, array, field lạ hoặc dữ liệu thừa phía sau object trả 400. Giá trị số có
giới hạn trên khớp kiểu cột (`NUMERIC(p,s)`, `SMALLINT`), vượt quá trả 422.

## Thêm bảng lookup mới

1. Thêm migration.
2. Thêm tên bảng vào `internal/modules/lookup/registry.go`.

`TestRegistryMatchesMigrations` sẽ fail nếu registry và migrations lệch nhau.

## Chạy

```bash
# Từ thư mục gốc repo (cần .env, xem infrastructure/database/README.md)
docker compose --profile app up -d --build
curl http://localhost:8081/healthz
```

Chạy local:

```bash
DB_PASSWORD=... go run ./cmd/api
```

Biến môi trường:

| Biến | Mặc định |
|---|---|
| `HTTP_PORT` | `8080` |
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` |
| `DB_NAME` | `product_db` |
| `DB_USER` | `product_user` |
| `DB_PASSWORD` | (bắt buộc) |
| `DB_SSLMODE` | `disable` |
| `DB_MAX_CONNS` | `10` |

Test: `go test ./...`
