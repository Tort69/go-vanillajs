export default class MovieRatingModal extends HTMLElement {
  // render(id) {
  //   const template = document.getElementById('movie-rating-modal-template')
  //   const content = template.content.cloneNode(true)
  //   this.appendChild(content)
  //   this._boundHandleMouseUp = this._handleCancel.bind(this)
  //   this.querySelector('.submit-btn').addEventListener('click', () =>
  //     app.saveToCollection(id, 'favorite')
  //   )
  //   const slider = document.getElementById('ratingSlider')
  //   const ratingValue = document.getElementById('ratingValue')
  //   slider.addEventListener('change', this._updateRating)
  // }
  // _updateRating() {
  //   debugger
  //   const value = parseInt(this.value)
  //   this.ratingValue.textContent = value
  // }

  connectedCallback() {
    const id = window.location.href.split('movies/')
    debugger
    this.render(id[1])
  }

  // open() {
  //   document.addEventListener('mouseup', this._boundHandleMouseUp)
  //   this.classList.add('modal--visible')
  // }

  // close() {
  //   this.classList.remove('modal--visible')
  //   document.removeEventListener('mouseup', this._boundHandleMouseUp)
  // }

  // _handleCancel(e) {
  //   e.preventDefault()
  //   if (e.target.matches('.cancel-btn')) {
  //     this.close()
  //   }
  //   if (!e.target.closest('.rating-container')) {
  //     this.close()
  //   }
  // }
}
customElements.define('movie-rating-modal', MovieRatingModal)
