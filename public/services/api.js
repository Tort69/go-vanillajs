import Router from './Router.js'

export const API = {
  baseURL: '/api/',
  getTopMovies: async () => {
    return await API.fetch('movies/top')
  },
  getRandomMovies: async () => {
    return await API.fetch('movies/random')
  },
  getMovieById: async (id) => {
    return await API.fetch(`movies/${id}`)
  },
  searchMovies: async (q, order, genre) => {
    return await API.fetch(`movies/search/`, { q, order, genre })
  },
  getGenres: async () => {
    return await API.fetch('genres/')
  },
  register: async (name, email, password) => {
    return await API.send('account/register/', { name, email, password })
  },
  authenticate: async (email, password) => {
    try {
      return await API.send('account/authenticate/', { email, password })
    } catch (e) {
      showError('Unable to authorize by credentials')
    }
  },
  verifyEmail: async (token) => {
    try {
      return await API.fetch(`account/verify/`, { token })
    } catch (e) {
      showError('Unable verify token')
    }
  },
  resendVerifyEmail: async (email) => {
    try {
      return await API.fetch(`account/resendVerifyEmail/`, { email })
    } catch (e) {
      showError('Unable send mail')
    }
  },
  deleteAccount: async () => {
    return await API.send('account/deleteAccount/')
  },
  getFavorites: async () => {
    try {
      return await API.fetch('account/favorites/')
    } catch (e) {
      app.Router.go('account/')
    }
  },
  getWatchlist: async () => {
    try {
      return await API.fetch('account/watchlist/')
    } catch (e) {
      app.Router.go('account/')
    }
  },
  saveToCollection: async (movie_id, collection) => {
    try {
      return await API.send('account/save-to-collection/', {
        movie_id,
        collection,
      })
    } catch (e) {
      app.Router.go('account/')
    }
  },
  deleteToCollection: async (movie_id, collection) => {
    try {
      const response = await API.send('account/delete-to-collection/', {
        movie_id,
        collection,
      })
      return response
    } catch (e) {
      app.showError('Internal Server Error', false)
    }
  },
  send: async (service, args) => {
    try {
      debugger
      const response = await fetch(API.baseURL + service, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: app.Store.jwt ? `Bearer ${app.Store.jwt}` : null,
        },
        body: JSON.stringify(args),
      })
      if (response.status == 201) {
        localStorage.setItem('unverifiedEmail', args.email)
        Router.go('/account/verifyEmail')
      }
      if (response.status == 403) {
        localStorage.setItem('unverifiedEmail', args.email)
        Router.go('/account/verifyEmail')
      }
      const result = await response.json()
      return result
    } catch (e) {
      app.showError('Internal Server Error', false)
    }
  },
  fetch: async (service, args) => {
    try {
      const queryString = args ? new URLSearchParams(args).toString() : ''
      const fullQueryString = args ? '?' + queryString : ''
      const response = await fetch(API.baseURL + service + fullQueryString, {
        headers: {
          Authorization: app.Store.jwt ? `Bearer ${app.Store.jwt}` : null,
        },
      })

      const result = await response.json()
      return result
    } catch (e) {
      app.showError('Internal Server Error', false)
    }
  },
}

export default API
