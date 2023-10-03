import { Application, Controller} from "@hotwired/stimulus"
const application = Application.start()

application.register("user--new-user", class extends Controller {
    static targets = [ "isRegistered", "registrationForm" ]

    connect() {
        if (this.isRegisteredTarget.value === "true") {
            location.href = "/"
        }
    }

    enableFields() {
        console.log("enableFields")
    }

    disableFields() {
        const form = this.registrationFormTarget
        for (const field of form.elements) {
            field.disabled = true
        }
        form.querySelector("button[type=submit]").innerHTML = "Registering..."
        form.querySelector("button[type=submit]").disabled = true
    }
})