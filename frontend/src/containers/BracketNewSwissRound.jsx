import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import New from './New';

// TODO: reuse stuff with BracketPrepare

const teamIdsSelector = createSelector(
  selectors.teamSeasons,
  (_, props) => props.season && props.season.get('id'),
  (x, seasonId) =>
    seasonId && x.toList()
      .filter(y => y.get('seasonId') === seasonId)
      .map(y => y.get('teamId')),
);

const teamsOptionsSelector = selectors.options(createSelector(
  selectors.teams,
  teamIdsSelector,
  (teams, teamIds) => teamIds && teams.toList().filter(
    x => teamIds.includes(x.get('id')),
  ),
));

@connect(createStructuredSelector({
  teamsOptions: teamsOptionsSelector,
}), actions)
export default class BracketNewSwissRound extends Component {
  static propTypes = {
    season: ImmutablePropTypes.map,
    bracket: ImmutablePropTypes.map,
    bracketPath: PropTypes.string,

    teamsOptions: ImmutablePropTypes.list,

    loadTeamSeasons: PropTypes.func.isRequired,
    patchBracket: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ season, teamsOptions }, prevProps) {
    if (!season) return;
    if (
      !prevProps ||
      season !== prevProps.season ||
      (!teamsOptions && prevProps.teamsOptions)
    ) {
      this.props.loadTeamSeasons(
        { seasonId: season.get('id') }, undefined, ['teamId'],
      );
    }
  }

  patcher = (...args) => {
    if (!this.props.bracket) return;
    this.props.patchBracket(this.props.bracket.get('id'), ...args);
  };

  // intentional force reload, because easy
  handleSuccess = () => { location.href = this.props.bracketPath; };

  render() {
    const { teamsOptions } = this.props;
    return (
      <New
        isAdminOnly
        resourceName="bracket"
        fields={{
          teams: {
            name: 'The teams',
            type: 'multi',
            options: teamsOptions ? teamsOptions.toJS() : [],
          },
          defaultTime: { name: 'Default match time', type: 'datetime' },
          reportMinutes: { name: 'Minutes for report submission', value: '0' },
          action: { type: 'hidden', value: 'new-swiss-round' },
        }}
        callForAction="generate new swiss round"
        createResource={this.patcher}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
