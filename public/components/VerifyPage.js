import API from '../services/API.js'

export default class VerifyPage extends HTMLElement {
  connectedCallback() {
    const template = document.getElementById('template-verify')
    const content = template.content.cloneNode(true)
    this.appendChild(content)
    const text = document.getElementById('verifyText')
    const email = window.localStorage.getItem('unverifiedEmail')
    text.textContent = `A confirmation link has been sent to ${email} mail address, please confirm your email address`
  }
}
customElements.define('verify-page', VerifyPage)
