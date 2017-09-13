import { PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { Link } from 'react-router';

import Component from '../utils/Component';

export default class TournamentSeasons extends Component {
  static propTypes = {
    tournamentSeasons: ImmutablePropTypes.list,
    tournamentPath: PropTypes.string,
  };

  render() {
    const { tournamentSeasons, tournamentPath } = this.props;
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          All seasons
          <ul>{tournamentSeasons.toJS().map(x => (
            <li key={x.id}>
              <Link to={`${tournamentPath}/${x.slug}`}>{x.name}</Link>
            </li>
          ))}</ul>
        </div>
      </div>
    );
  }
}
