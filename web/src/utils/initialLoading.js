export const removeInitialLoading = () => {
  if (typeof window !== 'undefined' && typeof window.__removeGvaLoading === 'function') {
    window.__removeGvaLoading()
    return
  }
  if (typeof document !== 'undefined') {
    document.getElementById('gva-loading-box')?.remove()
  }
}
