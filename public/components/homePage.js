import API from '../services/api.js'
import MovieItemComponent from './MovieItem.js'

export default class HomePage extends HTMLElement {
  async render() {
    const template = document.getElementById('template-home')
    const content = template.content.cloneNode(true)
    this.appendChild(content)

    const topMovies = await API.getTopMovies()
    this.renderMoviesInList(topMovies, this.querySelector('#top-10 ul'))

    const randomMovies = await API.getRandomMovies()
    this.renderMoviesInList(randomMovies, this.querySelector('#random ul'))
  }

  renderMoviesInList(movies, ul) {
    ul.innerHTML = ''
    movies.forEach((movie) => {
      const li = document.createElement('li')
      li.appendChild(new MovieItemComponent(movie))
      ul.appendChild(li)
    })
  }
  connectedCallback() {
    this.render()
  }
}
customElements.define('home-page', HomePage)
