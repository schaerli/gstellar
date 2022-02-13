(function() {
  const dropButtons = document.querySelectorAll(".snaphost-drop-button")
  const alert = document.querySelector("#successfulRestore")

  dropButtons.forEach(dropButton =>
    dropButton.addEventListener("click", (event) => {
      dropButton.disabled = true
      const sure = window.confirm("Are you sure?")

      if(sure) {
        const spinner = event.target.nextElementSibling
        spinner.style.display = "inline-block"

        fetch("/snapshots/drop?snapshot_id=" + event.target.dataset.snapshotId, {
          headers: {
            'Content-Type': 'application/json'
          }
        })
        .then(function (response) {
          response.json().then(data => {
            alert.innerHTML = data["Message"]
            alert.style.display = "inline-block"
            window.location = '/'
          })
        })
        .then(function (body) {
          dropButton.disabled = false
          spinner.style.display = "none"
        });
        } else {
          dropButton.disabled = false
          spinner.style.display = "none"
        }
    })
  )
})();