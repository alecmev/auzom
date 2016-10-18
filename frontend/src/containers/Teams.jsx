import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';
import { Link } from 'react-router';

import * as actions from '../actions';
import * as selectors from '../selectors';
import Component from '../utils/Component';

const activeTeamsSelector = createSelector(
  selectors.teams,
  teams =>
    teams.toList().filter(x =>
      x.get('disbandedAt') === null
    ).sort((a, b) =>
      a.get('name').localeCompare(b.get('name'))
    ),
);

@connect(createStructuredSelector({
  activeTeams: activeTeamsSelector,
}), actions)
export default class Teams extends Component {
  static propTypes = {
    activeTeams: ImmutablePropTypes.list,

    loadTeams: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.props.loadTeams(
      { disbandedAt: '%00' },
      this.props.activeTeams.reduce((r, x) => r.concat(x.get('id')), []),
    );
  }

  render() {
    const { activeTeams } = this.props;
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          All active teams
          <span className={this.cn({ d: 'secondary' })}>
            &nbsp;[{activeTeams.size}]
          </span>
          <ul>{activeTeams.toJS().map(team =>
            <li key={team.id}>
              <Link to={`/teams/${team.id}`}>
                [{team.abbr}] {team.name}
              </Link>
            </li>
          )}</ul>
        </div>
      </div>
    );
  }
}
