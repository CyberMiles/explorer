import { observable, action, runInAction } from "mobx"
import { BlockAPI } from "../api"

export class BlockStore {
  @observable isLoading = false
  @observable error = undefined
  @observable blocksRegistry = observable.map()

  getBlock(height) {
    return this.blocksRegistry.get(height)
  }

  @action
  loadBlock(height, { acceptCached = false } = {}) {
    this.error = undefined
    if (acceptCached) {
      const block = this.getBlock(height)
      if (block) return Promise.resolve(block)
    }
    this.isLoading = true
    return BlockAPI.get(height)
      .then(
        block => {
          runInAction(() => {
            this.blocksRegistry.set(height, block)
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

export default new BlockStore()
