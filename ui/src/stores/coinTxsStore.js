import { observable, action, runInAction } from "mobx"
import { AccountAPI } from "../api"

export class CoinTxsStore {
  @observable isLoading = false
  @observable error = undefined
  @observable address = undefined
  @observable coinTxs = []

  @action
  loadCoinTxs(address) {
    this.error = undefined
    this.isLoading = true
    this.coinTxs.clear()
    return AccountAPI.getCoinTxs(address)
      .then(
        txs => {
          runInAction(() => {
            this.coinTxs = txs
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

export default new CoinTxsStore()
