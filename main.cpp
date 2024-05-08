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

void openToWrite(std::ostream &file, uint16_t textToWrite) {
  int8_t parts[2];

  parts[0] = (textToWrite) & 0xFF;
  parts[1] = (textToWrite >> 8) & 0xFF;

  file.write((char *)parts, 2);
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

void fileInit(std::ostream &file, uint8_t *sizeOfWord) {
  const std::string id = "pice";
  const uint8_t version = 1;
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

int main() {
  uint16_t write[2] = {0xFFFF, 0xFEFF};
  uint8_t sizeOfWord = 2;

  std::string name = "1";

  std::fstream edit;

  uint8_t dataSize = 0;
  uint8_t curDataSize;

// code for specific open
  if(!isExist("bin/" + name + ".bin")) {
    edit.open("bin/" + name + ".bin",std::ios_base::binary | std::fstream::out);
    fileInit(edit, &sizeOfWord);
    curDataSize = 0;
  }
  else {
    edit.open("bin/" + name + ".bin",std::ios::binary | std::fstream::in | std::fstream::out | std::ios::ate);
    edit.seekg(5, std::ios::beg);
    edit.read((char*) &curDataSize, 1);
    edit.seekp(0, std::ios::end);
  }
// ------------------------------------
  for (char i = 0; i < 2; i++) {
    openToWrite(edit, write[i]);
    dataSize++;
  }
  
  curDataSize += dataSize;

  edit.seekp(5, std::ios::beg);
  edit.write((char *)&curDataSize, 1);

  edit.close();
  // In hex editor file read from right to left

  std::fstream outFile;

  outFile.open("bin/" + name + ".bin", std::ios::binary | std::ios::in);

  outFile.close();

  return 0;
}