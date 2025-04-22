import API from '../services/API.js'

export default class VerifyPage extends HTMLElement {
  token = null

  async render(token) {
    try {
      debugger
      const response = await API.verifyEmail(token)
    } catch (e) {
      app.showError(e, false)
      return
    }
    const template = document.getElementById('template-verify')
    const content = template.content.cloneNode(true)
    this.appendChild(content)
    if (response) {
      const htmlP = document.getElementById('verifyText')
      htmlP.innerText = 'Mail confirmed successful'
    } else {
      const htmlP = document.getElementById('verifyText')
      htmlP.innerText = 'Mail unconfirmed'
    }
  }

  connectedCallback() {
    const token = this.params[0]

    this.render(token)
  }
}
customElements.define('verif-mail-page', VerifyPage)
