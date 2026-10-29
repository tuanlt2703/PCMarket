let loaded = false

// adminlte.js wires its widgets on load; Vue renders the layout later, so load it after mount
// and re-run its lifecycle (it supports Turbo-style re-init) when the layout is mounted again.
export async function initAdminLte() {
  if (!loaded) {
    await import('@adminlte/js/adminlte.esm.js')
    loaded = true
    return
  }
  document.dispatchEvent(new Event('turbo:before-render'))
  document.dispatchEvent(new Event('turbo:load'))
}
