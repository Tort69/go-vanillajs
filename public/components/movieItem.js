export default class MovieItemComponent extends HTMLElement {
  constructor(movie, isDelete, title) {
    super()
    this.movie = movie
    this.isDelete = isDelete || false
    this.title = title || false
  }

  connectedCallback() {
    console.log(this.isDelete, this.title)
    const url = '/movies/' + this.movie.id

    if (this.isDelete) {
      this.innerHTML = `
                  <div class="movieArticle">
                  <a onclick="app.Router.go('${url}')">
                      <article>
                        <img src="${this.movie.poster_url}" alt="${this.movie.title} Poster" onerror="this.onerror=null; this.src='/images/image-not-found.jpg'">
                        <p>${this.movie.title} (${this.movie.release_year})</p>
                      </article>
                      </a>
                    <div class="movieArrow"> </div>
                  </div>
              `
      switch (this.title) {
        case 'Movie Watchlist':
          document
            .querySelector('.movieArrow')
            .addEventListener('click', (e) => {
              e.preventDefault()
              app.deleteToCollection(this.movie.id, 'watchlist')
            })
          break
        case 'Favorite Movies':
          document
            .querySelector('.movieArrow')
            .addEventListener('click', (e) => {
              e.preventDefault()
              app.deleteToCollection(this.movie.id, 'favorite')
            })
          break
      }
    } else {
      this.innerHTML = `
                    <a onclick="app.Router.go('${url}')">
                        <article>
                            <img src="${this.movie.poster_url}" alt="${this.movie.title} Poster"
                            onerror="this.onerror=null; this.src='/images/image-not-found.jpg'">
                            <p>${this.movie.title} (${this.movie.release_year})</p>
                        </article>
                    </a>
                `
    }
  }
}

customElements.define('movie-item', MovieItemComponent)
