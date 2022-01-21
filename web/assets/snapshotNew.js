(function() {
  const submitButton = document.querySelector("input[type='submit']");
  const spinner = document.querySelector(".gstellar-spinner");

  submitButton.addEventListener("click", (event) => {
    submitButton.disabled = true
    spinner.style.display = "inline-block"
    submitButton.form.submit()
  })

})();