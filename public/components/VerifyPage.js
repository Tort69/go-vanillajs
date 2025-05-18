import API from '../services/API.js'

export default class VerifyPage extends HTMLElement {
  connectedCallback() {
    const template = document.getElementById('template-verify')
    const content = template.content.cloneNode(true)
    this.appendChild(content)
  }
}
customElements.define('verify-page', VerifyPage)
