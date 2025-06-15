import API from '../services/API.js'
import MovieItemComponent from './MovieItem.js'

export default class MovieDetailsPage extends HTMLElement {
  response = null

  async render(id) {
    try {
      this.response = await API.getMovieById(id)
      this
    } catch (e) {
      app.showError()
      return
    }
    const template = document.getElementById('template-movie-details')
    const content = template.content.cloneNode(true)
    this.appendChild(content)
    const movieListStatus = this.response.movie.status
    if (movieListStatus.includes('In Favorite')) {
      this.querySelector('#btnFavorites').textContent = 'Unlist Favorite'
      this.querySelector('#btnFavorites').addEventListener('click', () => {
        app.deleteToCollection(this.response.movie.id, 'favorite')
      })
    } else {
      this.querySelector('#btnFavorites').addEventListener('click', () => {
        app.saveToCollection(this.response.movie.id, 'favorite')
      })
    }

    if (movieListStatus.includes('In Watchlist')) {
      this.querySelector('#btnWatchlist').textContent = 'Unlist Watchlist'
      this.querySelector('#btnWatchlist').addEventListener('click', () => {
        app.deleteToCollection(this.response.movie.id, 'watchlist')
      })
    } else {
      this.querySelector('#btnWatchlist').addEventListener('click', () => {
        app.saveToCollection(this.response.movie.id, 'watchlist')
      })
    }

    this.querySelector('h2').textContent = this.response.movie.title
    this.querySelector('h3').textContent = this.response.movie.tagline
    this.querySelector('img').src = this.response.movie.poster_url
    this.querySelector('#trailer').dataset.url = this.response.movie.trailer_url
    this.querySelector('#overview').textContent = this.response.movie.overview
    this.querySelector('#metadata').innerHTML = `
            <dt>Release Date</dt>
            <dd>${this.response.movie.release_year}</dd>
            <dt>Score</dt>
            <dd>${this.response.movie.score} / 10</dd>
            <dt>Original languae</dt>
            <dd>${this.response.movie.language}</dd>
        `

    const ulGenres = this.querySelector('#genres')
    ulGenres.innerHTML = ''
    this.response.movie.genres.forEach((genre) => {
      const li = document.createElement('li')
      li.textContent = genre.name
      ulGenres.appendChild(li)
    })

    const ulRelatedMovies = this.querySelector('#related-movies')
    ulRelatedMovies.textContent = ''
    this.response.related_movies.forEach((movie) => {
      const li = document.createElement('li')
      li.appendChild(new MovieItemComponent(movie))
      ulRelatedMovies.appendChild(li)
    })

    const ulCast = this.querySelector('#cast')
    ulCast.innerHTML = ''
    this.response.movie.casting.forEach((actor) => {
      const li = document.createElement('li')
      li.innerHTML = `
                <img src="${
                  actor.image_url ?? '/images/generic_actor.jpg'
                }" onerror="this.onerror=null; this.src='/images/generic_actor.jpg'" alt="Picture of ${
        actor.last_name
      }" >
                <p>${actor.first_name} ${actor.last_name}</p>
            `
      ulCast.appendChild(li)
    })
  }

  connectedCallback() {
    const id = this.params[0]

    this.render(id)
  }
}
customElements.define('movie-details-page', MovieDetailsPage)
