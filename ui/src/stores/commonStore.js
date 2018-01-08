import { observable, action, runInAction } from "mobx"
import { StatusAPI } from "../api"

class CommonStore {
  @observable appName = "Explorer"
  @observable appLoaded = false
  @observable isLoading = false
  @observable status = {}
  @observable error = undefined

  @action
  loadStatus() {
    this.isLoading = true
    return StatusAPI.get()
      .then(
        status => {
          runInAction(() => {
            this.status = status
          })
        },
        error => {
          runInAction(() => {
            this.error = error.message
          })
        }
      )
      .finally(
        action(() => {
          this.isLoading = false
        })
      )
  }

  @action
  setAppLoaded() {
    this.appLoaded = true
  }
}

export default new CommonStore()
