# Lunar Lander Monte Carlo Tree Search

## Action Space

There are four discrete actions available:

- 0: Do nothing
- 1: Fire left orientation engine
- 2: Fire main engine
- 3: Fire right orientation engine

## Observation Space

The state is an 8-dimensional vector consisting of:

1. The coordinates of the lander in x & y
2. Its linear velocities in x & y
3. Its angle
4. Its angular velocity
5. Two booleans representing whether each leg is in contact with the ground

## Rewards

After every step, a reward is granted. The total reward of an episode is the sum of the rewards for all the steps within that episode.

### Reward Details

- Proximity to the landing pad: Reward increases as the lander gets closer to the landing pad and decreases as it moves further away.
- Speed: Reward increases as the lander slows down and decreases as it moves faster.
- Angle: Reward decreases the more the lander is tilted (angle not horizontal).
- Leg Contact: Reward increases by 10 points for each leg in contact with the ground.
- Engine Usage:
  - Side engine: Reward decreases by 0.03 points for each frame a side engine is firing.
  - Main engine: Reward decreases by 0.3 points for each frame the main engine is firing.
- Episode Outcome:
  - Crashing: -100 points
  - Landing safely: +100 points

### Solution Criteria

An episode is considered a solution if it scores at least 200 points.

## Starting State

The lander starts at the top center of the viewport with a random initial force applied to its center of mass.

## Episode Termination

The episode finishes if:

- The lander crashes (the lander body gets in contact with the moon)
- The lander gets outside of the viewport (x coordinate is greater than 1)
