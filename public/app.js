import homePage from './components/homePage.js'

window.addEventListener('DOMContentLoaded', (e) => {
  document.querySelector('main').appendChild(new homePage())
})

window.app = {
  search: (event) => {
    event.preventDefault()
    const keywords = document.querySelector('input[type=search]').value
  },
}
