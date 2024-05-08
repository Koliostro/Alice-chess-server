#include <cstdint>
#include <cstdio>
#include <cstdlib>
#include <sys/stat.h>
#include <fstream>
#include <ios>
#include <iostream>

void initDesk() {

  for (char y = 0; y <= 7; y++) {
    for (char x = 0; x <= 7; x++) {
      std::cout << ". ";
    }

    std::cout << "  ";

    for (char x = 0; x <= 7; x++) {
      std::cout << ". ";
    }

    std::cout << std::endl;
  }
}

uint16_t openReadFile(std::istream &file) {
  {
    uint16_t value;
    uint8_t parts[2];

    // read 2 bytes from the file
    file.read((char *)parts, 2);
    // construct the 16-bit value from those bytes

    value = (parts[1] << 8) + parts[0];

    return value;
  }
}

// Writing part

void openToWrite(std::ostream &file, uint16_t textToWrite) {
  int8_t parts[2];

  parts[0] = (textToWrite) & 0xFF;
  parts[1] = (textToWrite >> 8) & 0xFF;

  file.write((char *)parts, 2);
}

void fileInit(std::ostream &file, uint8_t *sizeOfWord, const uint8_t version) {
  const std::string id = "pice";
  const uint8_t emptyData = 0;
  char * sizeWrittable = (char *) sizeOfWord;

  for (char i = 0; i < 4; i++) {
    file.write(&id[i], 1);
  }

  file.write((char *) &version, 1);
  file.write((char *) &emptyData, 1);
  file.write(sizeWrittable, 1);
}

bool isExist(const std::string &filepath) {
  struct stat buf;
   if (stat(filepath.c_str(), &buf) != -1) {
        return true;
    }
    return false;
}

uint8_t preapereFile(std::string filePath, std::fstream &file, uint8_t sizeOfWord, const uint8_t version) {
  uint8_t dataSize = 0;
  uint8_t curDataSize;

// code for specific open
  if(!isExist(filePath)) {
    file.open(filePath,std::ios_base::binary | std::fstream::out);
    fileInit(file, &sizeOfWord, version);
    curDataSize = 0;
  }
  else {
    file.open(filePath,std::ios::binary | std::fstream::in | std::fstream::out | std::ios::ate);
    file.seekg(5, std::ios::beg);
    file.read((char*) &curDataSize, 1);
    file.seekp(0, std::ios::end);
  }
// ------------------------------------  
  return curDataSize;
}

void endOfedit(std::fstream &file, uint8_t dataSize) {
  file.seekp(5, std::ios::beg);
  file.write((char *)&dataSize, 1);

  file.close();
}

// end of writing part
/*
  TODOS:

  create some classes for read and write and maybe in another file.

*/

int main() {
  uint16_t write[2] = {0xFFFF, 0xFEFF};
  uint8_t sizeOfWord = 2;
  const uint8_t ver = 1;

  std::string name = "1";
  std::string filePath = "bin/1.bin";

  std::fstream edit;

  uint8_t dataSize = preapereFile(filePath, edit, sizeOfWord, ver);

  // write here

  for (char i = 0; i < 2; i++) {
    openToWrite(edit, write[i]);
    dataSize++;
  }

  // ----------------

  endOfedit(edit, dataSize);
  
  // In hex editor file read from right to left

  std::fstream outFile;

  outFile.open("bin/" + name + ".bin", std::ios::binary | std::ios::in);

  outFile.close();

  return 0;
}