const Store = {
  jwt: null,
  get loggedIn() {
    if (this.jwt == 'null') {
      return false
    }
    return this.jwt !== null
  },
}

if (localStorage.getItem('jwt')) {
  Store.jwt = localStorage.getItem('jwt')
}

const proxiedStore = new Proxy(Store, {
  set: (target, prop, value) => {
    switch (prop) {
      case 'jwt':
        target[prop] = value
        console.log(value)
        localStorage.setItem('jwt', value)
        break
    }
    return true
  },
})

export default proxiedStore
