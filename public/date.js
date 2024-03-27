function formatDate(date) {
    let dateFormat = new Intl.DateTimeFormat('en', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    })
    return dateFormat.format(Date.parse(date))
}

// Set today's date
function updateDate() {
    let todayDate = new Date()
    let timeField = document.getElementById("today-date")
    timeField.innerText = formatDate(todayDate.toISOString().split('T')[0])

}

// Starts interval to update today's date every hour at 00:00
function callEveryHour() {
    setInterval(updateDate, 1000 * 60 * 60);
}

// Begin above interval
function beginDateUpdater() {
    let nextDate = new Date();
    nextDate.setHours(nextDate.getHours() + 1);
    nextDate.setMinutes(0);
    nextDate.setSeconds(0);

    const difference = nextDate - new Date();
    setTimeout(callEveryHour, difference);
}