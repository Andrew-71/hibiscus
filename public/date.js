// Format time in "Jan 02, 2006" format
function formatDate(date) {
    let dateFormat = new Intl.DateTimeFormat([langName, "en"], {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    })
    return dateFormat.format(date)
}

// Set today's date
function updateDate() {
    let timeField = document.getElementById("today-date")
    if (graceActive) {
        let graceField = document.getElementById("grace")
        graceField.hidden = false
    }
    timeField.innerText = formatDate(Date.now())

}

// Start interval to update today's date every hour at 00:00
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