import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import * as selectors from '../selectors';
import Component from '../utils/Component';

@connect(createStructuredSelector({
  teams: selectors.teams,
}))
export default class BracketStanding extends Component {
  static propTypes = {
    standing: ImmutablePropTypes.map.isRequired,
    odd: PropTypes.bool,
    hlTeamId: PropTypes.string,
    place: PropTypes.number,
    placeTo: PropTypes.number,
    playoffs: PropTypes.bool,
    swiss: PropTypes.bool,
    aceSwiss: PropTypes.bool,
    seasonPath: PropTypes.string,

    onMouseEnter: PropTypes.func,
    onMouseLeave: PropTypes.func,

    teams: ImmutablePropTypes.map,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleMouseEnter = () =>
    this.props.onMouseEnter &&
    this.props.onMouseEnter(this.props.standing.get('teamId'));
  handleMouseLeave = () =>
    this.props.onMouseLeave &&
    this.props.onMouseLeave();
  handleClick = () => this.context.router.push(
    `${this.props.seasonPath}/participants/${this.props.standing.get('teamId')}`
  );

  render() {
    const {
      standing, odd, hlTeamId, place, placeTo, playoffs, swiss, aceSwiss, teams,
    } = this.props;

    const s = standing.toJS();
    let name = s.teamId;
    const team = teams.get(s.teamId);
    if (team) name = team.get('name');

    let mod;
    if (place === 1) mod = 'gold';
    else if (place === 2) mod = 'silver';
    else if (place === 3) mod = 'bronze';

    const events = {
      onMouseEnter: this.handleMouseEnter,
      onMouseLeave: this.handleMouseLeave,
      onClick: this.handleClick,
    };

    return (
      <tr className={this.cni({ m: [
        (s.teamId === hlTeamId) && 'highlight',
        odd && 'odd',
      ] })}>
        {place &&
          <td
            className={this.cn({ d: 'place' })}
            rowSpan={placeTo && ((placeTo - place) + 1)}
          >
            {place}{placeTo && ` - ${placeTo}`}
          </td>
        }{place &&
          <td
            className={this.cn({ d: 'divider', m: mod })}
            rowSpan={placeTo && ((placeTo - place) + 1)}
          />
        }
        <td className={this.cn({ d: 'team' })} {...events}>{name}</td>
        {swiss && <td {...events}>{s.byes || ''}</td>}
        {!playoffs && <td {...events}>{s.scoreWon}</td>}
        {aceSwiss && <td {...events}>{s.medianBuchholz}</td>}
        {!playoffs && <td {...events}>{s.mapsWon}</td>}
        {!playoffs && <td {...events}>{s.roundsWon}</td>}
        {!playoffs && <td {...events}>{s.rawScoreRatio === '+Inf' ?
          'N/A' : (+s.rawScoreRatio).toFixed(4)}</td>}
        {!playoffs && <td {...events}>{s.rawScoreWon}</td>}
        {!playoffs && <td {...events}>{s.rawScoreLost}</td>}
        {playoffs && <td {...events}>{s.matchLossWeight}</td>}
        {playoffs && <td {...events}>{s.matchesLost}</td>}
      </tr>
    );
  }
}
