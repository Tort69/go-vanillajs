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
                        <p class="movieScore">${this.movie.score}</p>
                      </article>
                      </a>
                <div class="movieArrow">
                  <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg"><g id="SVGRepo_bgCarrier" stroke-width="0"></g><g id="SVGRepo_tracerCarrier" stroke-linecap="round" stroke-linejoin="round"></g><g id="SVGRepo_iconCarrier"> <path d="M16 8L8 16M8.00001 8L16 16" stroke="#000000" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"></path> </g></svg>
                </div>
                  </div>
              `
      switch (this.title) {
        case 'Movie Watchlist':
          document.querySelectorAll('.movieArrow').forEach((element) =>
            element.addEventListener('click', (e) => {
              e.preventDefault()
              app.deleteToCollection(this.movie.id, 'watchlist')
            })
          )
          break
        case 'Favorite Movies':
          document.querySelectorAll('.movieArrow').forEach((element) =>
            element.addEventListener('click', (e) => {
              e.preventDefault()
              app.deleteToCollection(this.movie.id, 'favorite')
            })
          )
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
