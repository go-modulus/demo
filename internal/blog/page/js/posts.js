import { Application, Controller} from "@hotwired/stimulus"
const application = Application.start()

application.register("blog--posts", class extends Controller {
    static targets = ['addPostForm']

    connect() {

    }

    openAddPopup() {
        console.log("enableFields")
    }
    disableFields() {
        console.log("disableFields")
    }
})