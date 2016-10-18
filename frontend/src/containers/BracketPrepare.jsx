import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import New from './New';

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
export default class BracketPrepare extends Component {
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
          mapsPerMatch: { name: 'Amount of maps per match' },
          maps: {
            type: 'hidden',
            value: [
              '12', '11', '4', '5', '9', '3', '6', '10',
              '1', '12', '10', '5', '11', '9',
            ], // TODO: ugh...
          },
          action: { type: 'hidden', value: 'prepare' },
        }}
        callForAction="prepare the bracket"
        createResource={this.patcher}
        onSuccess={this.handleSuccess}
      />
    );
  }
}
