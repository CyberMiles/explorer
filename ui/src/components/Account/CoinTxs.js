import React, { Component } from "react"
import { Header, Grid, Segment } from "semantic-ui-react"
import { inject, observer } from "mobx-react"
import { Link } from "react-router-dom"

import Coins from "./Coins"
import LoadingSpinner from "../common/LoadingSpinner"
import RedError from "../common/RedError"

class CoinTx extends Component {
  render() {
    const { inputs, outputs } = this.props.tx
    return (
      <Segment>
        <Header as="h4">Height: {this.props.height}</Header>
        <Grid stackable columns="equal">
          <Grid.Column>
            {inputs.map((input, index) => (
              <Input key={index} address={input.sender} coins={input.coins} />
            ))}
          </Grid.Column>
          <Grid.Column width={1}>==></Grid.Column>
          <Grid.Column>
            {outputs.map((output, index) => (
              <Input key={index} address={output.receiver} coins={output.coins} />
            ))}
          </Grid.Column>
        </Grid>
      </Segment>
    )
  }
}

@inject("accountStore")
@observer
class Input extends Component {
  render() {
    const { address, coins } = this.props
    const DivAddress =
      this.props.accountStore.address === address ? (
        <Segment>{address}</Segment>
      ) : (
        <Segment as={Link} to={"/account/" + address}>
          {address}
        </Segment>
      )
    return (
      <Segment.Group horizontal size="mini">
        {DivAddress}
        <Coins coins={coins} />
      </Segment.Group>
    )
  }
}

@inject("coinTxsStore")
@observer
class CoinTxs extends Component {
  componentDidMount() {
    const address = this.props.address
    this.props.coinTxsStore.loadCoinTxs(address)
  }
  componentWillReceiveProps(nextProps) {
    var addr_old = this.props.address
    var addr_new = nextProps.address
    if (addr_old !== addr_new) {
      this.props.coinTxsStore.loadCoinTxs(addr_new)
    }
  }
  render() {
    const { isLoading, error } = this.props.coinTxsStore
    if (isLoading) return <LoadingSpinner />
    if (error) return <RedError message={error} />

    const txs = this.props.coinTxsStore.coinTxs || []
    if (txs.length === 0) return <div />
    return (
      <div>
        <Header>Coin Transactions</Header>
        {txs.map((tx, index) => <CoinTx key={index} tx={tx.tx} height={tx.height} />)}
      </div>
    )
  }
}

export default CoinTxs
