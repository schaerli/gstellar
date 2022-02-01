(function() {
  const restoreButton = document.querySelector("#snaphostRestoreLink")
  const snapshotId = restoreButton.dataset.snapshotId

  restoreButton.addEventListener("click", (event) => {
    restoreButton.disabled = true
    const sure = window.confirm("Are you sure?")

    if(sure) {
      fetch("/get-data", opts).then(function (response) {
        return response.json();
      })
      .then(function (body) {
        //doSomething with body;
      });

    } else {
      restoreButton.disabled = false
    }

    // spinner.style.display = "inline-block"
    // submitButton.form.submit()
  })

})();