import API from '../services/API.js'
import MovieItemComponent from './MovieItem.js'

export default class ActorDetailPage extends HTMLElement {
  response = null

  async render(id) {
    try {
      this.response = await API.getActorDetail(id)
    } catch (e) {
      app.showError()
      return
    }
    const template = document.getElementById('template-actor-details')
    const content = template.content.cloneNode(true)
    this.appendChild(content)
    const movieListStatus = this.response.related_movies

    this.querySelector('h2').textContent =
      this.response.actor.first_name + ' ' + this.response.actor.last_name

    if (this.response.actor.image_url) {
      this.querySelector('img').src = this.response.actor.image_url
    } else this.querySelector('img').src = '/images/generic_actor.jpg'

    const ulRelatedMovies = this.querySelector('#actor-movies')
    ulRelatedMovies.textContent = ''
    this.response.related_movies.forEach((movie) => {
      const li = document.createElement('li')
      li.appendChild(new MovieItemComponent(movie))
      ulRelatedMovies.appendChild(li)
    })
  }

  connectedCallback() {
    const id = this.params[0]

    this.render(id)
  }
}
customElements.define('movie-actor-page', ActorDetailPage)
