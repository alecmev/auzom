import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';

import * as actions from '../actions';

import New from './New';

@connect(undefined, actions)
export default class SeasonApply extends Component {
  static propTypes = {
    seasonPath: PropTypes.string,
    season: ImmutablePropTypes.map,
    leaderships: ImmutablePropTypes.list,
    teams: ImmutablePropTypes.map,
    createTeamSeasonRequest: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleSuccess = x => this.context.router.push(
    `${this.props.seasonPath}/participants/${x.teamId}`
  );

  render() {
    const { season, leaderships, teams } = this.props;
    if (!season || !leaderships) return null;
    return (
      <New
        resourceName="application"
        fields={{
          teamId: {
            name: 'Your team',
            type: 'select',
            options: leaderships.toJS().map((x) => {
              const team = teams.get(x.teamId);
              return {
                label: (team && team.get('name')) || x.teamId,
                value: x.teamId,
              };
            }),
          },
          seasonId: { type: 'hidden', value: season.get('id') },
        }}
        callForAction="submit season application"
        createResource={this.props.createTeamSeasonRequest}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
