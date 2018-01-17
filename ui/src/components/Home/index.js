import React, { Component } from "react"
import { Grid } from "semantic-ui-react"

import RecentBlocks from "./RecentBlocks"
import RecentCoinTxs from "./RecentCoinTxs"
import RecentStakeTxs from "./RecentStakeTxs"

export default class Home extends Component {
  render() {
    return (
      <Grid columns={1}>
        <Grid.Column style={{ paddingTop: 0 }}>
          <RecentBlocks />
        </Grid.Column>
      </Grid>
    )
  }
}
