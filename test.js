import { default as bar, default as foo }from'./engine';

import imgMarioRunning from './sprites/marioRunning.png';
import imgMarioIdle from './sprites/marioIdle.png';
import imgMarioJumping from './sprites/marioJumping.png';
import imgGroundBlock from './sprites/groundBlock.png';
import imgCoin from './sprites/coin.png';
import imgGoomba from './sprites/goomba.png';
import imgQuestionBlock from './sprites/questionBlock.png';

import audioMainTheme from './sounds/mainTheme.mp3';
import audioJump from './sounds/jump.wav';
import audioCoin from './sounds/coin.wav'; 
import audioStomp from './sounds/stomp.wav';

const {
  addEntityType,
  addEntity,
  removeEntity,
  createMap,

  loadSprite,
  drawSprite,

  loadSound,
  playSound,

  settings,
  keys,
  time,
  keyCode,
  camera,

  moveAndCheckForObstacles,
  checkCollision,
} = engine;

// загрузка спрайтов будет не в юзеркоде наверн
const sprMarioRunning = loadSprite(imgMarioRunning, -8, -16, 2);
const sprMarioJumping = loadSprite(imgMarioJumping, -8, -16);
const sprMarioIdle = loadSprite(imgMarioIdle, -8, -16);
const sprGroundBlock = loadSprite(imgGroundBlock, 0, 0);
const sprGoomba = loadSprite(imgGoomba, -8, -16, 2);
const sprCoin = loadSprite(imgCoin, 0, 0, 2);
const sprQuestionBlock = loadSprite(imgQuestionBlock, 0, 0, 1);

const sndMainTheme = loadSound(audioMainTheme);
const sndJump = loadSound(audioJump);
const sndCoin = loadSound(audioCoin);
const sndStomp = loadSound(audioStomp);
//playSound(sndMainTheme, true);


function updateWall(wall) {
  drawSprite(sprGroundBlock, wall);
}

const ENTITY_TYPE_WALL = addEntityType('#', updateWall, {
  bbox: {
    left: 0,
    top: 0,
    width: settings.tileSize,
    height: settings.tileSize
  }
});


let playerStartX = 0;
let playerStartY = 0;

function updateMario(mario) {
  if (!mario.isInitialized) {
    playerStartX = mario.x;
    playerStartY = mario.y;

    mario.isInitialized = true;
  }

  const absSpeedX = Math.abs(mario.speedX);
  const dir = mario.direction || 1;

  if (!mario.isOnGround) {
    drawSprite(sprMarioJumping, mario, 0, dir);
  } else if (absSpeedX > 1) {
    drawSprite(
      sprMarioRunning,
      mario,
      0.03 * absSpeedX,
      dir
    );
  } else {
    drawSprite(sprMarioIdle, mario, 0, dir);
  }

  const keyLeft = keys[keyCode.ARROW_LEFT];
  const keyRight = keys[keyCode.ARROW_RIGHT];
  const keySpace = keys[keyCode.SPACE];

  const metersPerSecondSq = 60;
  const friction = 10;

  const accelX = (keyRight.isDown - keyLeft.isDown) * metersPerSecondSq;

  if (keyRight.isDown) {
    mario.direction = 1;
  }
  if (keyLeft.isDown) {
    mario.direction = -1;
  }

  mario.speedX += accelX * time.deltaTime;
  mario.speedY += settings.gravity * time.deltaTime;
  mario.speedX *= 1 - friction * time.deltaTime;

  const { vertWall } = moveAndCheckForObstacles(mario, [ENTITY_TYPE_WALL, ENTITY_TYPE_QUESTION_BLOCK]);
  mario.isOnGround = vertWall && vertWall.y <= mario.y;

  if (keySpace.wentDown && mario.isOnGround) {
    mario.speedY = -12;
    playSound(sndJump);
  }

  // question blocks
  if (vertWall && vertWall.type === ENTITY_TYPE_QUESTION_BLOCK && vertWall.y < mario.y) {
    const coin = addEntity(ENTITY_TYPE_COIN);
    coin.x = vertWall.x;
    coin.y = vertWall.y - settings.tileSize;
    coin.isFlying = true;
  }

  const hitEnemy = checkCollision(mario, [ENTITY_TYPE_GOOMBA]);
  if (hitEnemy) {
    if (hitEnemy.y > mario.y) {
      removeEntity(hitEnemy);
      mario.speedY = -15;
      playSound(sndStomp);
    } else {
      removeEntity(mario);
      // TODO: сделать меню после проигрыша
      const newMario = addEntity(ENTITY_TYPE_MARIO);
      newMario.x = playerStartX;
      newMario.y = playerStartY;
    }
  }

  if (mario.y > 30) {
    removeEntity(mario);
    const newMario = addEntity(ENTITY_TYPE_MARIO);
    newMario.x = playerStartX;
    newMario.y = playerStartY;
  }

  const hitCoin = checkCollision(mario, [ENTITY_TYPE_COIN]);
  if (hitCoin) {
    removeEntity(hitCoin);
    playSound(sndCoin);
  }

  camera.x = mario.x;
  camera.y = 6;
}

const ENTITY_TYPE_MARIO = addEntityType('@', updateMario, {
  bbox: {
    left: -0.45,
    top: -1,
    width: 0.9,
    height: 1
  },
});


function updateGoomba(goomba) {
  const { horizWall } = moveAndCheckForObstacles(goomba, [ENTITY_TYPE_WALL, ENTITY_TYPE_QUESTION_BLOCK]);
  if (horizWall) {
    if (goomba.x < horizWall.x) {
      goomba.speedX = -2;
    } else {
      goomba.speedX = 2;
    }
  }
  goomba.speedY += settings.gravity * time.deltaTime;
  drawSprite(sprGoomba, goomba, 3 * time.deltaTime);
}

const ENTITY_TYPE_GOOMBA = addEntityType('G', updateGoomba, {
  bbox: {
    left: -0.5,
    top: -1,
    width: 1,
    height: 1
  },
  speedX: 2
});


function updateCoin(coin) {
  if (!coin.isInitialized) {
    coin.startY = coin.y;
    coin.isInitialized = true;
  }

  drawSprite(sprCoin, coin, 2 * time.deltaTime);
  if (coin.isFlying) {
    coin.y -= 4 * time.deltaTime;
    if (coin.y < coin.startY - 1) {
      removeEntity(coin);
    }
  }
}

const ENTITY_TYPE_COIN = addEntityType('0', updateCoin, {
  bbox: {
    left: 0,
    top: 0,
    width: 1,
    height: 1,
  }
});

function updateQuestionBlock(block) {
  drawSprite(sprQuestionBlock, block);
}

const ENTITY_TYPE_QUESTION_BLOCK = addEntityType('?', updateQuestionBlock, {
  bbox: {
    left: 0,
    top: 0,
    width: 1,
    height: 1,
  }
});

const asciiMapRows = [
  ' # ##########                                             ',
  '                                                          ',
  '      ###      0000                                       ',
  '##  #####      ####                                       ',
  '#                                                         ',
  '        #  G        GGGGGGGGGGGGGGGGGGG                   ',
  '#         ###                                             ',
  '#                                                         ',
  '  ?                                                       ',
  '#                                                         ',
  '#   @    #   G G   #        0000000000000000000           ',
  '######   #################################################'
];

createMap(asciiMapRows);
 