function getRemainingTime(time) {
  const RESEND_DELAY = 60
  if (time <= 0) return 0
  const elapsed = Math.floor((Date.now() - time) / 1000)
  return Math.max(RESEND_DELAY - elapsed, 0)
}

export default function startTimer() {
  this.button = document.getElementById('resend-button')
  this.timerDisplay = document.getElementById('timer')
  const lastEmailSentTime =
    window.localStorage.getItem('lastEmailSentTime') || 0

  seconds = getRemainingTime(lastEmailSentTime)
  let remaining = seconds
  this.button.disabled = true

  intervalId = setInterval(() => {
    remaining--

    if (remaining <= 0) {
      clearInterval(intervalId)
      this.button.disabled = false
      this.timerDisplay.textContent = ''
      return
    }
    this.timerDisplay.textContent = `Повторная отправка через ${remaining} сек.`
  }, 1000)

  this.timerDisplay.textContent = `Повторная отправка через ${remaining} сек.`
  this.button.disabled = true
}
