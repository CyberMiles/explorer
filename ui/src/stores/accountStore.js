import { observable, action, runInAction } from "mobx"
import { AccountAPI } from "../api"

export class AccountStore {
  @observable isLoading = false
  @observable error = undefined
  @observable accountsRegistry = observable.map()

  getAccount(address) {
    return this.accountsRegistry.get(address)
  }

  @action
  loadAccount(address, { acceptCached = false } = {}) {
    this.error = undefined
    if (acceptCached) {
      const account = this.getAccount(address)
      if (account) return Promise.resolve(account)
    }
    this.isLoading = true
    return AccountAPI.get(address)
      .then(
        account => {
          runInAction(() => {
            this.accountsRegistry.set(address, account)
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
}

export default new AccountStore()
