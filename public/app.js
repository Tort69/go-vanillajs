import YouTubeEmbed from './components/YouTubeEmbed.js'
import Router from './services/Router.js'
import MovieItemComponent from './components/MovieItem.js'
import HomePage from './components/HomePage.js'
import MovieDetailsPage from './components/MovieDetailPage.js'
import MoviesPage from './components/MoviesPage.js'
import { routes } from './services/Routes.js'
import ResetPasswordModal from './components/ResetPasswordModal.js'
import Store from './services/Store.js'
import API from './services/api.js'
import startTimer, { getRemainingTime } from './utils/startTimer.js'
import MovieRatingModal from './components/MovieRatingModal.js'
import PaginationComponent from './components/PaginationComponent.js'

window.app = {
  API,
  Router,
  Store,
  showError: (
    message = 'There was an error loading the page',
    goToHome = true
  ) => {
    document.querySelector('#alert-modal p').textContent = message
    document.querySelector('#alert-modal').showModal()
    if (goToHome) app.Router.go('/')
    return
  },
  closeError: () => {
    document.getElementById('alert-modal').close()
  },
  openModal: () => {
    document.querySelector('main').appendChild(new ResetPasswordModal())
    document.querySelector('password-modal').open()
  },
  closeModal: () => {
    document.querySelector('password-modal').close()
  },
  openMovieRatingModal: () => {
    const modal = document.querySelector('#rating-modal')
    modal.showModal()
    const slider = document.getElementById('ratingSlider')
    const ratingValue = document.getElementById('ratingValue')
    slider.addEventListener('change', () => {
      ratingValue.textContent = slider.value
    })
    document
      .querySelector('#rating-action-btn')
      .addEventListener('click', () => {
        const id = window.location.href.split('movies/')
        app.saveCollection(parseInt(id[1]), 'favorite', parseInt(slider.value))
        modal.close()
      })
  },
  closeMovieRatingModal: () => {
    document.getElementById('rating-modal').close()
  },
  search: (event) => {
    event.preventDefault()
    const keywords = document.querySelector('input[type=search]').value
    if (keywords.length > 1) {
      app.Router.go(`/movies?q=${keywords}`)
    }
  },
  searchOrderChange: (order) => {
    const urlParams = new URLSearchParams(window.location.search)
    const q = urlParams.get('q')
    const genre = urlParams.get('genre') ?? ''
    app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`)
  },
  searchFilterChange: (genre) => {
    const urlParams = new URLSearchParams(window.location.search)
    const q = urlParams.get('q')
    const order = urlParams.get('order') ?? ''
    app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`)
  },
  register: async (event) => {
    event.preventDefault()
    let errors = []
    const name = document.getElementById('register-name').value
    const email = document.getElementById('register-email').value
    const password = document.getElementById('register-password').value
    const passwordConfirm = document.getElementById(
      'register-password-confirm'
    ).value

    if (name.length < 4) errors.push('Enter your complete name')
    if (email.length < 8) errors.push('Enter your complete email')
    if (password.length < 6) errors.push('Enter a password with 6 characters')
    if (password != passwordConfirm) errors.push("Passwords don't match")
    if (errors.length == 0) {
      const response = await API.register(name, email, password)
      if (response.success) {
        localStorage.setItem('unverifiedEmail', email)
        app.Router.go('/account/verifyEmail')
      } else {
        app.showError(response.message, false)
      }
    } else {
      app.showError(errors.join('. '), false)
    }
  },
  handlerResetPassword: async (event) => {
    event.preventDefault()
    const formData = new FormData(event.target)
    const data = {
      currentPassword: formData.get('currentPassword'),
      newPassword: formData.get('newPassword'),
      newPasswordConfirm: formData.get('newPasswordConfirm'),
    }
    if (data.newPassword != data.newPasswordConfirm) {
      app.showError("Passwords don't match", false)
    }
    try {
      const response = await API.resetPassword(
        data.currentPassword,
        data.newPassword
      )
      if (response?.success) {
        app.closeModal()
        app.showError('Password successfully changed', false)
      }
    } catch (e) {
      app.showError(e.message, false)
    }
  },
  login: async (event) => {
    event.preventDefault()
    let errors = []
    const email = document.getElementById('login-email').value
    const password = document.getElementById('login-password').value

    if (email.length < 8) errors.push('Enter your complete email')
    if (password.length < 6) errors.push('Enter a password with 6 characters')
    if (errors.length == 0) {
      const response = await API.authenticate(email, password)
      if (!response) {
        return
      }
      if (response.success) {
        app.Store.jwt = response.jwt
        app.Router.go('/account/')
      }
    } else {
      app.showError(errors.join('. '), false)
    }
  },
  handlerResendVerifyEmail: async (event) => {
    event.preventDefault()
    const email = window.localStorage.getItem('unverifiedEmail')
    const seconds = localStorage.getItem('lastEmailSentTime')
    if (getRemainingTime(seconds) == 0) {
      try {
        localStorage.setItem('lastEmailSentTime', Date.now())
        const response = await API.resendVerifyEmail(email)
        if (!response) return
      } catch (e) {
        app.showError('Unable send mail', false)
        return
      }
    }
    startTimer()
  },
  logout: () => {
    localStorage.removeItem('jwt')
    app.Store.jwt = null
    app.Router.go('/')
  },
  deleteAccount: async () => {
    const response = await API.deleteAccount()
    if (response.success) {
      localStorage.removeItem('jwt')
      app.Store.jwt = null
      app.Router.go('/')
    }
  },
  saveCollection: async (movie_id, collection, score = null) => {
    if (app.Store.loggedIn) {
      try {
        const response = await API.saveToCollection(movie_id, collection, score)
        if (response.success) {
          switch (collection) {
            case 'favorite':
              app.Router.go('/account/favorites')
              break
            case 'watchlist':
              app.Router.go('/account/watchlist')
          }
        } else {
          app.showError("We couldn't save the movie.")
        }
      } catch (e) {
        console.log(e)
      }
    } else {
      app.Router.go('/account/')
    }
  },
  deleteToCollection: async (movie_id, collection) => {
    if (app.Store.loggedIn) {
      try {
        const response = await API.deleteToCollection(movie_id, collection)
        if (response.success) {
          switch (collection) {
            case 'favorite':
              app.Router.go('/account/favorites')
              break
            case 'watchlist':
              app.Router.go('/account/watchlist')
              break
          }
        } else {
          app.showError("We couldn't delete the movie.")
        }
      } catch (e) {
        console.log(e)
      }
    } else {
      app.Router.go('/account/')
    }
  },
  actorDetail: async (id) => {},
}
window.addEventListener('DOMContentLoaded', () => {
  app.Router.init()
})
