import ImmutablePropTypes from 'react-immutable-proptypes';

import Component from '../utils/Component';

export default class MatchDetails extends Component {
  static propTypes = {
    match: ImmutablePropTypes.map,
    matchReport: ImmutablePropTypes.map,
    matchRounds: ImmutablePropTypes.list,
    matchWidePenalties: ImmutablePropTypes.list,
    roundWidePenalties: ImmutablePropTypes.map,

    gameMaps: ImmutablePropTypes.map,
  };

  render() {
    const {
      match, matchReport, matchRounds, gameMaps,
      matchWidePenalties, roundWidePenalties,
    } = this.props;

    if (!match || !matchReport) return null;

    const report = matchReport.toJS();
    const isMatchDrawRaw = report.rawScoreX === report.rawScoreY;
    const isMatchDraw = report.scoreX === report.scoreY;
    const isOverriddenGlobally = (
      report.rawScoreXOverride !== null || report.rawScoreYOverride !== null
    );
    let prevMapId;
    let sideX;
    let sideY;
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ d: 'inner' })}>
          {matchRounds.toJS().map((x) => {
            let mapName;
            if (x.gameMapId !== prevMapId) {
              prevMapId = x.gameMapId;
              const map = gameMaps.get(x.gameMapId);
              if (map) {
                [mapName, sideX, sideY] = [
                  map.get('abbr'),
                  map.get('sideXAbbr'),
                  map.get('sideYAbbr'),
                ];
              }
            }

            const isDraw = x.rawScoreX === x.rawScoreY;
            const penalties = roundWidePenalties.get(x.id);
            const isOverridden = (x.rawScoreXOverride || x.rawScoreYOverride);
            return (
              <div key={x.id} className={this.cn({ d: 'round' })}>
                {mapName && <div className={this.cn({ d: 'map' })}>
                  {mapName}
                </div>}
                <div className={this.cn({
                  d: 'side',
                  m: 'left',
                })}>{x.isTeamXOnSideY ? sideY : sideX}</div>
                <div className={this.cn({
                  d: 'side',
                  m: 'right',
                })}>{x.isTeamXOnSideY ? sideX : sideY}</div>
                {x.isNotPlayed ?
                  <div className={this.cn({
                    d: 'notPlayed',
                    m: (isOverridden || isOverriddenGlobally) &&
                      'cancelledOut',
                  })}>
                    not played
                  </div> :
                  <div className={this.cn({ d: 'scores' })}>
                    <div className={this.cn({
                      d: 'score',
                      m: [
                        'left',
                        (isDraw && 'draw') ||
                          (x.rawScoreX > x.rawScoreY ? 'win' : 'loss'),
                        (x.rawScoreXOverride || report.rawScoreXOverride) &&
                          'cancelledOut',
                      ],
                    })}>{x.rawScoreX}</div>
                    <div className={this.cn({
                      d: 'score',
                      m: [
                        'right',
                        (isDraw && 'draw') ||
                          (x.rawScoreX < x.rawScoreY ? 'win' : 'loss'),
                        (x.rawScoreYOverride || report.rawScoreYOverride) &&
                          'cancelledOut',
                      ],
                    })}>{x.rawScoreY}</div>
                  </div>
                }
                {penalties && penalties.toJS().map(y => (
                  <div
                    key={y.id}
                    className={this.cn({ d: 'penalty' })}
                  >
                    <div className={this.cn({
                      d: 'score',
                      m: ['left', 'penalty',
                        (x.rawScoreXOverride || report.rawScoreXOverride) &&
                          'cancelledOut'],
                    })}>{-y.rawScoreX}</div>
                    <div className={this.cn({
                      d: 'score',
                      m: ['right', 'penalty',
                        (x.rawScoreYOverride || report.rawScoreYOverride) &&
                          'cancelledOut'],
                    })}>{-y.rawScoreY}</div>
                    <div className={this.cn({ d: 'reason' })}>
                      <span className={this.cn({ d: 'reasonRed' })}>
                        Penalty
                      </span> {y.reason}
                    </div>
                  </div>
                ))}
                {isOverridden &&
                  <div className={this.cn({ d: 'penalty' })}>
                    <div className={this.cn({
                      d: 'score',
                      m: ['left', 'override',
                        isOverriddenGlobally && 'cancelledOut'],
                    })}>{x.rawScoreXOverride}</div>
                    <div className={this.cn({
                      d: 'score',
                      m: ['right', 'override',
                        isOverriddenGlobally && 'cancelledOut'],
                    })}>{x.rawScoreYOverride}</div>
                    <div className={this.cn({ d: 'reason' })}>
                      <span className={this.cn({ d: 'reasonGreen' })}>
                        {x.isPenalOverride && 'Penal '}Override
                      </span> {x.overrideReason}
                    </div>
                  </div>
                }
              </div>
            );
          })}
          <div className={this.cn({ d: 'rule' })} />
          {matchWidePenalties && matchWidePenalties.toJS().map(x => (
            <div
              key={x.id}
              className={this.cn({ d: 'penalty' })}
            >
              <div className={this.cn({
                d: 'score',
                m: ['left', 'penalty',
                  report.rawScoreXOverride && 'cancelledOut'],
              })}>{-x.rawScoreX}</div>
              <div className={this.cn({
                d: 'score',
                m: ['right', 'penalty',
                  report.rawScoreYOverride && 'cancelledOut'],
              })}>{-x.rawScoreY}</div>
              <div className={this.cn({ d: 'reason' })}>
                <span className={this.cn({ d: 'reasonRed' })}>
                  Penalty
                </span> {x.reason}
              </div>
            </div>
          ))}
          {isOverriddenGlobally &&
            <div className={this.cn({ d: 'penalty' })}>
              <div className={this.cn({
                d: 'score',
                m: ['left', 'override'],
              })}>{report.rawScoreXOverride}</div>
              <div className={this.cn({
                d: 'score',
                m: ['right', 'override'],
              })}>{report.rawScoreYOverride}</div>
              <div className={this.cn({ d: 'reason' })}>
                <span className={this.cn({ d: 'reasonGreen' })}>
                  {report.isPenalOverride && 'Penal '}Override
                </span> {report.overrideReason}
              </div>
            </div>
          }
          <div className={this.cn({ d: 'scores' })}>
            <div
              className={this.cn({ d: 'map' })}
              title="tickets / kills / etc."
            >
              raw total
            </div>
            <div className={this.cn({
              d: 'score',
              m: ['left', (isMatchDrawRaw && 'draw') ||
                (report.rawScoreX > report.rawScoreY ? 'win' : 'loss'),
              ],
            })}>{report.rawScoreX}</div>
            <div className={this.cn({
              d: 'score',
              m: ['right', (isMatchDrawRaw && 'draw') ||
                (report.rawScoreX < report.rawScoreY ? 'win' : 'loss'),
              ],
            })}>{report.rawScoreY}</div>
          </div>
          <div className={this.cn({ d: 'rule' })} />
          {matchWidePenalties && matchWidePenalties.toJS().map(x => (
            <div
              key={x.id}
              className={this.cn({ d: 'penalty' })}
            >
              <div className={this.cn({
                d: 'score',
                m: ['left', 'penalty',
                  report.scoreXOverride && 'cancelledOut'],
              })}>{-x.scoreX}</div>
              <div className={this.cn({
                d: 'score',
                m: ['right', 'penalty',
                  report.scoreYOverride && 'cancelledOut'],
              })}>{-x.scoreY}</div>
              <div className={this.cn({ d: 'reason' })}>
                <span className={this.cn({ d: 'reasonRed' })}>
                  Penalty
                </span> {x.reason}
              </div>
            </div>
          ))}
          {(report.scoreXOverride !== null || report.scoreYOverride !== null) &&
            <div className={this.cn({ d: 'penalty' })}>
              <div className={this.cn({
                d: 'score',
                m: ['left', 'override'],
              })}>{report.scoreXOverride}</div>
              <div className={this.cn({
                d: 'score',
                m: ['right', 'override'],
              })}>{report.scoreYOverride}</div>
              <div className={this.cn({ d: 'reason' })}>
                <span className={this.cn({ d: 'reasonGreen' })}>
                  {report.isPenalOverride && 'Penal '}Override
                </span> {report.overrideReason}
              </div>
            </div>
          }
          <div className={this.cn({ d: 'scores' })}>
            <div className={this.cn({ d: 'map' })}>points</div>
            <div className={this.cn({
              d: 'score',
              m: ['left', (isMatchDraw && 'draw') ||
                (report.scoreX > report.scoreY ? 'win' : 'loss'),
              ],
            })}>{report.scoreX}</div>
            <div className={this.cn({
              d: 'score',
              m: ['right', (isMatchDraw && 'draw') ||
                (report.scoreX < report.scoreY ? 'win' : 'loss'),
              ],
            })}>{report.scoreY}</div>
          </div>
        </div>
      </div>
    );
  }
}
