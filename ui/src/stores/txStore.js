import { observable, action, runInAction } from "mobx"
import { TxAPI } from "../api"

export class TxStore {
  @observable isLoading = false
  @observable error = undefined
  @observable blocksRegistry = observable.map()

  getTx(txhash) {
    return this.blocksRegistry.get(txhash)
  }

  @action
  loadTx(txhash, { acceptCached = false } = {}) {
    this.error = undefined
    if (acceptCached) {
      const block = this.getTx(txhash)
      if (block) return Promise.resolve(block)
    }
    this.isLoading = true
    return TxAPI.get(txhash)
      .then(
        block => {
          runInAction(() => {
            this.blocksRegistry.set(txhash, block)
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

export default new TxStore()
