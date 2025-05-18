import API from '../services/API.js'

export default class ConfirmedMailPage extends HTMLElement {
  token = null

  async render(token) {
    try {
      const response = await API.verifyEmail(token)
      if (response) {
        app.Router.go('/')
      } else {
        app.showError(
          'Неудалось верифицировато почту, попробуйте другие данные',
          false
        )
        app.Router.go('/account/verifyEmail')
      }
    } catch (e) {
      app.showError(e, false)
      return
    }
  }

  connectedCallback() {
    const urlParams = new URLSearchParams(window.location.search)
    const token = urlParams.get('token')
    this.render(token)
  }
}
customElements.define('confirmed-mail-page', ConfirmedMailPage)
