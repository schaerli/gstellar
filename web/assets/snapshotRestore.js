(function() {
  const restoreButtons = document.querySelectorAll(".snaphost-restore-button")

  restoreButtons.forEach(restoreButton =>
    restoreButton.addEventListener("click", (event) => {
      restoreButton.disabled = true
      const sure = window.confirm("Are you sure?")

      if(sure) {
        const spinner = event.target.nextElementSibling
        spinner.style.display = "inline-block"

        fetch("/snapshots/restore?snapshot_id=" + event.target.dataset.snapshotId).then(function (response) {
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