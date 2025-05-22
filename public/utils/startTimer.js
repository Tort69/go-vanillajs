export function getRemainingTime(time) {
  const RESEND_DELAY = 55
  if (time <= 0) return 0
  const elapsed = Math.floor((Date.now() - time) / 1000)
  return Math.max(RESEND_DELAY - elapsed, 0)
}

export default function startTimer() {
  const button = document.getElementById('resend-button')
  const timerDisplay = document.getElementById('timer')
  const lastEmailSentTime =
    window.localStorage.getItem('lastEmailSentTime') || 0

  const seconds = getRemainingTime(lastEmailSentTime)
  let remaining = seconds
  button.disabled = true

  const intervalId = setInterval(() => {
    remaining--

    if (remaining <= 0) {
      clearInterval(intervalId)
      button.disabled = false
      timerDisplay.textContent = 'Resend'
      return
    }
    timerDisplay.textContent = `Повторная отправка через ${remaining} сек.`
  }, 1000)
  timerDisplay.textContent = `Повторная отправка через ${remaining} сек.`
  button.disabled = true
}
