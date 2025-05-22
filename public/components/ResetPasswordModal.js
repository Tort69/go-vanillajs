export default class PasswordModal extends HTMLElement {
  constructor() {
    super()
    this.className = 'password-modal'

    const template = document.getElementById('password-modal-template')
    const content = template.content.cloneNode(true)
    this.appendChild(content)
    // this._handleCancel = this._handleCancel.bind(this)
    this._boundHandleMouseUp = this._handleCancel.bind(this)
  }

  connectedCallback() {}

  open() {
    document.addEventListener('mouseup', this._boundHandleMouseUp)
    this.classList.add('password-modal--visible')
  }

  close() {
    this.classList.remove('password-modal--visible')
    this.querySelector('.password-modal__form').reset()
    document.removeEventListener('mouseup', this._boundHandleMouseUp)
  }

  _handleCancel(e) {
    e.preventDefault()
    if (e.target.matches('.password-modal__cancel-btn')) {
      this.close()
    }
    if (!e.target.closest('.password-modal__content')) {
      this.close()
    }
  }
}
customElements.define('password-modal', PasswordModal)
