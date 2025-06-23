export default class PaginationComponent extends HTMLElement {
  constructor() {
    super()
    this.currentPage = 1
    this.totalPages = 1
    this.pageSize = 50
    this.response = {}
  }

  connectedCallback() {
    this.render()
    this.addEventListeners()
  }

  render() {
    this.innerHTML = `
          <style>
              .pagination { display: flex; gap: 5px; margin-top: 20px; }
              button { padding: 5px 10px; cursor: pointer; }
              button.active { background: #007bff; color: white; }
          </style>
          <div class="pagination"></div>
      `
    this.updateButtons()
  }

  updateButtons() {
    const paginationDiv = this.querySelector('.pagination')
    paginationDiv.textContent = ''

    // Кнопка "Назад"
    const prevButton = document.createElement('button')
    prevButton.textContent = '←'
    prevButton.disabled = this.currentPage === 1
    prevButton.dataset.page = this.currentPage - 1
    paginationDiv.appendChild(prevButton)

    // Номера страниц
    const rangePage = this.totalPages / this.currentPage
    if (rangePage < 10) {
      for (let i = this.currentPage; i <= this.totalPages; i++) {
        const btn = document.createElement('button')
        btn.textContent = i
        btn.classList.toggle('active', i === this.currentPage)
        btn.dataset.page = i
        paginationDiv.appendChild(btn)
        btn.addEventListener('click', () => {
          const urlParams = new URLSearchParams(window.location.search)
          urlParams.set('page', i)
          urlParams.set('pageSize', this.pageSize)
          app.Router.go(window.location.pathname + window.location.search)
        })
      }
    } else {
      for (let i = this.currentPage; i <= this.currentPage + 9; i++) {
        if (this.currentPage + i <= this.currentPage + 7) {
          const btn = document.createElement('button')
          btn.textContent = i
          btn.classList.toggle('active', i === this.currentPage)
          btn.dataset.page = i
          paginationDiv.appendChild(btn)
          btn.addEventListener('click', () => {
            const urlParams = new URLSearchParams(window.location.search)
            urlParams.set('page', i)
            urlParams.set('pageSize', this.pageSize)
            app.Router.go(window.location.pathname + window.location.search)
          })
        } else {
          const btn = document.createElement('button')
          btn.textContent = this.totalPages - i
          btn.classList.toggle('active', i === this.currentPage)
          btn.dataset.page = i
          paginationDiv.appendChild(btn)
          btn.addEventListener('click', () => {
            const urlParams = new URLSearchParams(window.location.search)
            urlParams.set('page', i)
            urlParams.set('pageSize', this.pageSize)
            app.Router.go(window.location.pathname + window.location.search)
          })
        }
      }
    }

    // Кнопка "Вперед"
    const nextButton = document.createElement('button')
    nextButton.textContent = '→'
    nextButton.disabled = this.currentPage === this.totalPages
    nextButton.dataset.page = this.currentPage + 1
    paginationDiv.appendChild(nextButton)
  }

  addEventListeners() {
    this.addEventListener('click', (e) => {
      if (e.target.tagName === 'BUTTON') {
        const newPage = parseInt(e.target.dataset.page)
        if (!isNaN(newPage)) {
          this.dispatchEvent(
            new CustomEvent('page-changed', {
              detail: newPage,
              bubbles: true,
              composed: true,
            })
          )
        }
      }
    })
  }

  setPages(currentPage, totalPages) {
    this.currentPage = currentPage
    this.totalPages = totalPages
    this.updateButtons()
  }
}

customElements.define('pagination-component', PaginationComponent)
