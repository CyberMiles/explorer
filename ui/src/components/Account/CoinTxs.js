import React, { Component } from "react"
import { Segment, Table } from "semantic-ui-react"
import { inject, observer } from "mobx-react"
import { Link } from "react-router-dom"

import Coins from "./Coins"
import RedError from "../common/RedError"

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
    if (error !== undefined) return <RedError message={error} />

    const txs = this.props.coinTxsStore.coinTxs || []
    return (
      <Segment basic loading={isLoading}>
        <Table basic="very" compact>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell width={1}>Height</Table.HeaderCell>
              <Table.HeaderCell width={5}>From</Table.HeaderCell>
              <Table.HeaderCell width={5}>To</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {txs.map((tx, index) => <CoinTx key={index} tx={tx.tx} height={tx.height} />)}
          </Table.Body>
        </Table>
      </Segment>
    )
  }
}

class CoinTx extends Component {
  render() {
    const height = this.props.height
    const { inputs, outputs } = this.props.tx
    return (
      <Table.Row verticalAlign="top">
        <Table.Cell>{height}</Table.Cell>
        <Table.Cell>
          {inputs.map((input, index) => (
            <InputOutput key={index} address={input.sender} coins={input.coins} />
          ))}
        </Table.Cell>
        <Table.Cell>
          {outputs.map((output, index) => (
            <InputOutput key={index} address={output.receiver} coins={output.coins} />
          ))}
        </Table.Cell>
      </Table.Row>
    )
  }
}

@inject("coinTxsStore")
class InputOutput extends Component {
  render() {
    const { address, coins } = this.props
    const DivAddress =
      this.props.coinTxsStore.address === address ? (
        <span>{address}</span>
      ) : (
        <Link to={"/account/" + address}>{address}</Link>
      )

    return (
      <Table basic="very" compact fixed singleLine>
        <Table.Body>
          <Table.Row verticalAlign="top">
            <Table.Cell width={3} style={{ paddingTop: "0px" }}>
              {DivAddress}
            </Table.Cell>
            <Table.Cell width={1} textAlign="right" style={{ paddingTop: "0px" }}>
              <Coins coins={coins} />
            </Table.Cell>
          </Table.Row>
        </Table.Body>
      </Table>
    )
  }
}

export default CoinTxs
