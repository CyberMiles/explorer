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
    this.address = address
    this.coinTxs.clear()
    return AccountAPI.getCoinTxs(address)
      .then(
        txs => {
          runInAction(() => {
            if (address === this.address) this.coinTxs = txs
          })
        },
        error => {
          runInAction(() => {
            if (address === this.address) this.error = error.message
          })
        }
      )
      .finally(
        action(() => {
          if (address === this.address) this.isLoading = false
        })
      )
  }
}

export default new CoinTxsStore()
