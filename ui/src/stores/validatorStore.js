import { observable, action, runInAction } from "mobx"
import { ValidatorAPI } from "../api"

export class ValidatorStore {
  @observable isLoading = false
  @observable error = undefined
  @observable height = undefined
  @observable validators = []

  @action
  loadValidators(height, { acceptCached = false } = {}) {
    this.error = undefined
    this.isLoading = true
    this.height = height
    this.validators.clear()
    return ValidatorAPI.get(height)
      .then(
        ret => {
          runInAction(() => {
            if (height === this.height) this.validators = ret.validators
          })
        },
        error => {
          runInAction(() => {
            if (height === this.height) this.error = error.message
          })
        }
      )
      .finally(
        action(() => {
          if (height === this.height) this.isLoading = false
        })
      )
  }
}

export default new ValidatorStore()
