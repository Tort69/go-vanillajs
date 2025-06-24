export default class PaginationComponent extends HTMLElement {
  constructor() {
    super()
    this.currentPage = 1
    this.totalPages = 1
    this.pageSize = 50
  }

  connectedCallback() {
    this.render()
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
    prevButton.addEventListener('click', () => {
      const urlParams = new URLSearchParams(window.location.search)
      const query = urlParams.get('q') || ''
      const order = urlParams.get('order') || ''
      const genre = urlParams.get('genre') || ''
      app.Router.go(
        `/movies?q=${query}&order=${order}&genre=${genre}&page=${prevButton.dataset.page}&pageSize=${this.pageSize}`
      )
    })
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
          const query = urlParams.get('q') || ''
          const order = urlParams.get('order') || ''
          const genre = urlParams.get('genre') || ''
          this.currentPage = i
          app.Router.go(
            `/movies?q=${query}&order=${order}&genre=${genre}&page=${i}&pageSize=${this.pageSize}`
          )
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
            const query = urlParams.get('q') || ''
            const order = urlParams.get('order') || ''
            const genre = urlParams.get('genre') || ''
            this.currentPage = i
            app.Router.go(
              `/movies?q=${query}&order=${order}&genre=${genre}&page=${i}&pageSize=${this.pageSize}`
            )
          })
        } else {
          const btn = document.createElement('button')
          btn.textContent = this.totalPages - i
          btn.classList.toggle('active', i === this.currentPage)
          btn.dataset.page = i
          paginationDiv.appendChild(btn)
          btn.addEventListener('click', () => {
            const urlParams = new URLSearchParams(window.location.search)
            const query = urlParams.get('q') || ''
            const order = urlParams.get('order') || ''
            const genre = urlParams.get('genre') || ''
            this.currentPage = i
            app.Router.go(
              `/movies?q=${query}&order=${order}&genre=${genre}&page=${i}&pageSize=${this.pageSize}`
            )
          })
        }
      }
    }

    // Кнопка "Вперед"
    const nextButton = document.createElement('button')
    nextButton.textContent = '→'
    nextButton.disabled = this.currentPage === this.totalPages
    nextButton.dataset.page = this.currentPage + 1
    nextButton.addEventListener('click', () => {
      const urlParams = new URLSearchParams(window.location.search)
      const query = urlParams.get('q') || ''
      const order = urlParams.get('order') || ''
      const genre = urlParams.get('genre') || ''
      app.Router.go(
        `/movies?q=${query}&order=${order}&genre=${genre}&page=${nextButton.dataset.page}&pageSize=${this.pageSize}`
      )
    })
    paginationDiv.appendChild(nextButton)
  }

  setPages(currentPage, totalPages) {
    this.currentPage = currentPage
    this.totalPages = totalPages
    this.render()
  }
}

customElements.define('pagination-component', PaginationComponent)
