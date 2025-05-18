import AccountPage from '../components/AccountPage.js'
import FavoritePage from '../components/FavoritePage.js'
import HomePage from '../components/HomePage.js'
import LoginPage from '../components/LoginPage.js'
import MovieDetailsPage from '../components/MovieDetailPage.js'
import MoviesPage from '../components/MoviesPage.js'
import RegisterPage from '../components/RegisterPage.js'
import VerifyPage from '../components/VerifyPage.js'
import WatchlistPage from '../components/WatchlistPage.js'
import ConfirmedMailPage from '../components/ConfirmedMailPage.js'

export const routes = [
  {
    path: '/',
    component: HomePage,
    loggedIn: false,
  },
  {
    path: '/movies',
    component: MoviesPage,
    loggedIn: false,
  },
  {
    path: /\/movies\/(\d+)/,
    component: MovieDetailsPage,
    loggedIn: false,
  },
  {
    path: '/account/register',
    component: RegisterPage,
    loggedIn: false,
  },
  {
    path: '/account/login',
    component: LoginPage,
    loggedIn: false,
  },
  {
    path: '/account/',
    component: AccountPage,
    loggedIn: true,
  },
  {
    path: '/account/favorites',
    component: FavoritePage,
    loggedIn: true,
  },
  {
    path: '/account/watchlist',
    component: WatchlistPage,
    loggedIn: true,
  },
  {
    path: '/account/verifyEmail',
    component: VerifyPage,
    loggedIn: false,
  },
  {
    path: '/account/verify',
    component: ConfirmedMailPage,
    loggedIn: false,
  },
]
