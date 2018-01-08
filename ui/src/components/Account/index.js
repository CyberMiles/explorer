import React, { Component } from "react"
import { inject, observer } from "mobx-react"
import { withRouter } from "react-router-dom"
import { Header, Grid } from "semantic-ui-react"

import LoadingSpinner from "../common/LoadingSpinner"
import RedError from "../common/RedError"
import Coins from "./Coins"
import CoinTxs from "./CoinTxs"

@inject("accountStore")
@withRouter
@observer
export default class Account extends Component {
  componentDidMount() {
    const address = this.props.match.params.address
    this.props.accountStore.loadAccount(address, { acceptCached: true })
  }
  componentWillReceiveProps(nextProps) {
    var addr_old = this.props.match.params.address
    var addr_new = nextProps.match.params.address
    if (addr_old !== addr_new) {
      this.props.accountStore.loadAccount(addr_new, { acceptCached: true })
    }
  }

  render() {
    const { isLoading, error } = this.props.accountStore
    if (isLoading) return <LoadingSpinner />
    if (error) return <RedError message={error} />

    const address = this.props.match.params.address
    const account = this.props.accountStore.getAccount(address)

    if (!account) return <RedError message="Can't load account" />
    return (
      <Grid>
        <Grid.Row columns={1}>
          <Grid.Column>
            <Header>Address</Header>
            <Header.Subheader>{address}</Header.Subheader>
          </Grid.Column>
        </Grid.Row>
        <Grid.Row columns={2}>
          <Grid.Column>
            <Coins title="Coins" coins={account.coins} />
          </Grid.Column>
          <Grid.Column>
            <Coins title="Credit" coins={account.credit} />
          </Grid.Column>
        </Grid.Row>
        <Grid.Row columns={1}>
          <Grid.Column>
            <CoinTxs address={address} />
          </Grid.Column>
        </Grid.Row>
      </Grid>
    )
  }
}
