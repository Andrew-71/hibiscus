// Format time in "Jan 02, 2006" format
function formatDate(date) {
    let dateFormat = new Intl.DateTimeFormat([langName, "en"], {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    })
    return dateFormat.format(date)
}

async function graceActive() {
    const response = await fetch("/api/grace");
    const active = await response.text();
    return active === "true"
}

// Set today's date and grace status
function updateTime() {
    document.getElementById("today-date").innerText = formatDate(Date.now());
    graceActive().then(v => document.getElementById("grace").hidden = !v)
}

// Start interval to update time at start of every minute
function callEveryMinute() {
    setInterval(updateTime, 1000 * 60);
}

// Begin above interval
function beginTimeUpdater() {
    const difference = (60 - new Date().getSeconds()) * 1000;
    setTimeout(callEveryMinute, difference);
    setTimeout(updateTime, difference);
    updateTime();
}

// This does NOT properly sanitize, and assumes a well-meaning user
function sanitize(title) {
    return title
        .trim()
        .replace(/ +/g, '-')
        .replace(/[!*'();:@&=+$,\/?#\[\]]/g, '')
}

// Open a new note
function newNote(text_prompt) {
    name = sanitize(prompt(text_prompt + ':'))
    window.location.replace('/notes/' + name)
}