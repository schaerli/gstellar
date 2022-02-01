(function() {
  const restoreButtons = document.querySelectorAll(".snaphost-restore-button")
  const alert = document.querySelector("#successfulRestore")

  restoreButtons.forEach(restoreButton =>
    restoreButton.addEventListener("click", (event) => {
      restoreButton.disabled = true
      const sure = window.confirm("Are you sure?")

      if(sure) {
        const spinner = event.target.nextElementSibling
        spinner.style.display = "inline-block"

        fetch("/snapshots/restore?snapshot_id=" + event.target.dataset.snapshotId, {
          headers: {
            'Content-Type': 'application/json'
          }
        })
        .then(function (response) {
          response.json().then(data => {
            alert.innerHTML = data["Message"]
            alert.style.display = "inline-block"
          })
        })
        .then(function (body) {
          restoreButton.disabled = false
          spinner.style.display = "none"
        });
        } else {
          restoreButton.disabled = false
          spinner.style.display = "none"
        }
    })
  )
})();