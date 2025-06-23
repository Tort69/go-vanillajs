import API from '../services/api.js'
import MovieItemComponent from './MovieItem.js'

export default class MoviesPage extends HTMLElement {
  constructor() {
    super()
    this.movies = []
    this.currentPage = 1
    this.pageSize = 10
    this.totalCount = 0
    this.isLoading = false
  }

  async loadGenres() {
    const genres = await API.getGenres()
    const select = this.querySelector('#filter')
    select.innerHTML = `
		<option value=''>Filter by Genre</option>
	`
    genres.forEach((genre) => {
      var option = document.createElement('option')
      option.value = genre.id
      option.textContent = genre.name
      select.appendChild(option)
    })
  }
  async render() {
    const urlParams = new URLSearchParams(window.location.search)
    const query = urlParams.get('q') ?? ''
    const order = urlParams.get('order') ?? ''
    const genre = urlParams.get('genre') ?? ''
    const releaseYear = urlParams.get('releaseYear') ?? ''
    let page = urlParams.get('page') ?? ''
    let pageSize = urlParams.get('pageSize') ?? ''
    if (query) {
      this.querySelector('h2').textContent = `'${query}' movies`
      this.render(query)
    }
    debugger
    // const movies = await fetch(
    //   '/api/movies/search/?query=&order=&genre=&releaseYear=&page=&pageSize='
    // )

    if (page == '') {
      page = 1
    }
    if (pageSize == '') {
      pageSize = 50
    }

    const movies = await API.searchMovies(
      query,
      order,
      genre,
      releaseYear,
      page,
      pageSize
    )
    // const movies = await response.Json()

    const totalPages = Math.ceil(movies.totalCount / movies.pageSize)
    const pagination = this.querySelector('pagination-component')
    pagination.setPages(this.currentPage, totalPages)
    this.querySelector('ul').append(pagination)

    const ulMovies = this.querySelector('ul')
    ulMovies.innerHTML = ''
    if (movies && movies.length > 0) {
      movies.forEach((movie) => {
        const li = document.createElement('li')
        li.appendChild(new MovieItemComponent(movie))
        ulMovies.appendChild(li)
      })
    } else {
      ulMovies.textContent = 'There are no movies with your search'
    }

    await this.loadGenres()

    if (order) this.querySelector('#order').value = order
    if (genre) this.querySelector('#filter').value = genre
  }

  connectedCallback() {
    const template = document.getElementById('template-movies')
    const content = template.content.cloneNode(true)
    this.appendChild(content)
    this.render()
  }
}
customElements.define('movies-page', MoviesPage)
